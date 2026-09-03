package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func projectFail(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "Error:") {
		return err
	}
	return fmt.Errorf("Error: %s", msg)
}

func runProject(home string, args []string) error {
	setCommand(currentCtx, "project", args)

	remaining, err := lessflags.
		Help("-h,--help", projectHelp()).
		HelpNoExit().
		StopOnFirstArg().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) == 0 {
		return projectFail(fmt.Errorf("tsk project: subcommand required"))
	}
	switch remaining[0] {
	case "add":
		return runProjectAdd(home, remaining[1:])
	case "tree":
		return runProjectTree(home, remaining[1:])
	case "list":
		return runProjectRegistryList(home, remaining[1:])
	case "which":
		return runProjectWhich(home, remaining[1:])
	case "register":
		return runProjectRegister(home, remaining[1:])
	case "unregister":
		return runProjectUnregister(home, remaining[1:])
	default:
		return projectFail(fmt.Errorf("tsk project: unknown subcommand %q", remaining[0]))
	}
}

func runProjectAdd(home string, args []string) error {
	var dirFlag, projectName string
	var notes []string
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		String("--project", &projectName).
		StringSlice("--note", &notes).
		Help("-h,--help", projectAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 1 {
		return projectFail(fmt.Errorf("tsk project add: title required"))
	}
	title := remaining[0]
	if title == "" {
		return projectFail(fmt.Errorf("tsk project add: title required"))
	}
	for i, text := range notes {
		notes[i] = strings.TrimSpace(text)
		if notes[i] == "" {
			return projectFail(fmt.Errorf("tsk project add: --note text required"))
		}
	}

	cwd, ref, err := resolveAddProject(home, dirFlag, projectName)
	if err != nil {
		return projectFail(err)
	}

	if err := storage.EnsureLayout(home); err != nil {
		return err
	}
	id, err := storage.NextID(home)
	if err != nil {
		return err
	}

	slug := storage.Slugify(title)
	stage := "create"
	now := storage.NowTimestamp(id)
	relPath := storage.InboxRelPath(id, title)
	taskDir := filepath.Join(home, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Join(taskDir, "context"), 0o755); err != nil {
		return err
	}

	task := storage.Task{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Labels:       []string{},
		TopicPath:    storage.NullTopicPath,
		Cwd:          cwd,
		Project:      &ref,
		Stage:        stage,
		CreatedAt:    now,
		UpdatedAt:    now,
		StageHistory: []storage.StageHistoryEntry{},
	}
	if err := storage.WriteTask(taskDir, task); err != nil {
		return err
	}
	if err := storage.WriteIndex(home, id, relPath); err != nil {
		return err
	}
	for _, text := range notes {
		existing, err := storage.ReadTopicNotes(taskDir)
		if err != nil {
			return projectFail(err)
		}
		note := storage.TopicNote{
			TS:   storage.NowTimestamp(len(existing) + 1),
			Text: text,
		}
		if err := storage.AppendTopicNote(taskDir, note); err != nil {
			return projectFail(fmt.Errorf("tsk project add: append note: %w", err))
		}
	}
	autoCwd := pathfmt.TildeHome(mainRepoDirForAuto(cwd))
	if err := storage.UpsertProjectAuto(home, ref, autoCwd); err != nil {
		return projectFail(fmt.Errorf("tsk project add: update projects-auto: %w", err))
	}
	fmt.Println(id)
	return nil
}

func resolveAddProject(home, dirFlag, projectName string) (cwd string, ref storage.ProjectRef, err error) {
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return "", storage.ProjectRef{}, err
	}

	projectName = strings.TrimSpace(projectName)
	if projectName != "" {
		entry, ok := storage.FindProjectByName(reg, projectName)
		if !ok {
			return "", storage.ProjectRef{}, fmt.Errorf("tsk project add: unknown project %q (see tsk project register --help)", projectName)
		}
		probe := strings.TrimSpace(dirFlag)
		if probe == "" && entry.Cwd != "" {
			probe = pathfmt.Expand(entry.Cwd)
		}
		if probe == "" {
			probe, err = resolveProbeDir("")
			if err != nil {
				return "", storage.ProjectRef{}, err
			}
		} else {
			probe, err = resolveProbeDir(probe)
			if err != nil {
				return "", storage.ProjectRef{}, err
			}
		}
		if origin, oerr := gitNormalizedOrigin(probe); oerr == nil && origin != "" {
			return probe, storage.ProjectRef{Origin: origin}, nil
		}
		if entry.Origin != "" {
			return probe, storage.ProjectRef{Origin: entry.Origin}, nil
		}
		if entry.Name == "" {
			return "", storage.ProjectRef{}, fmt.Errorf("tsk project add: registered project %q has no name", projectName)
		}
		return probe, storage.ProjectRef{Name: entry.Name}, nil
	}

	probeDir, err := resolveProbeDir(dirFlag)
	if err != nil {
		return "", storage.ProjectRef{}, err
	}
	if origin, oerr := gitNormalizedOrigin(probeDir); oerr == nil && origin != "" {
		return probeDir, storage.ProjectRef{Origin: origin}, nil
	}
	if entry, ok := findRegistryByAbsCwd(reg, probeDir); ok && entry.Name != "" {
		return probeDir, storage.ProjectRef{Name: entry.Name}, nil
	}
	return "", storage.ProjectRef{}, fmt.Errorf("tsk project add: no git origin and not registered (see tsk project register --help)")
}

func runProjectWhich(home string, args []string) error {
	var dirFlag string
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		Help("-h,--help", projectWhichHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project which: unexpected arguments"))
	}
	probeDir, err := resolveProbeDir(dirFlag)
	if err != nil {
		return projectFail(err)
	}
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return projectFail(err)
	}

	origin, _ := gitNormalizedOrigin(probeDir)
	entry, hasReg := findRegistryByAbsCwd(reg, probeDir)
	if origin != "" {
		if e, ok := storage.FindProjectByOrigin(reg, origin); ok {
			entry, hasReg = e, true
		}
	}

	if origin == "" && !hasReg {
		return projectFail(fmt.Errorf("tsk project which: no git origin and not registered (see tsk project register --help)"))
	}

	if origin != "" {
		fmt.Printf("origin:  %s\n", origin)
	}
	if hasReg && entry.Name != "" {
		fmt.Printf("name:    %s\n", entry.Name)
	} else if origin != "" {
		_, base, _ := NormalizeOriginURL("https://" + origin)
		if base == "" {
			base = filepath.Base(origin)
		}
		fmt.Printf("name:    %s\n", base)
	}
	fmt.Printf("cwd:     %s\n", probeDir)
	return nil
}

func runProjectRegister(home string, args []string) error {
	var name, cwdFlag, originFlag string
	remaining, err := lessflags.
		String("--name", &name).
		String("--cwd", &cwdFlag).
		String("--origin", &originFlag).
		Help("-h,--help", projectRegisterHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project register: unexpected arguments"))
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return projectFail(fmt.Errorf("tsk project register: --name required"))
	}

	probeDir, err := resolveProbeDir(cwdFlag)
	if err != nil {
		return projectFail(err)
	}
	cwdTilde := pathfmt.TildeHome(probeDir)

	origin := strings.TrimSpace(originFlag)
	if origin == "" {
		if o, oerr := gitNormalizedOrigin(probeDir); oerr == nil {
			origin = o
		}
	} else {
		// Allow raw URLs; normalize when possible.
		if id, _, nerr := NormalizeOriginURL(origin); nerr == nil {
			origin = id
		}
	}

	reg, err := storage.ReadProjects(home)
	if err != nil {
		return projectFail(err)
	}
	if _, ok := storage.FindProjectByName(reg, name); ok {
		return projectFail(fmt.Errorf("tsk project register: name %q already registered", name))
	}

	entry := storage.ProjectEntry{
		Name:   name,
		Cwd:    cwdTilde,
		Origin: origin,
	}
	reg.Projects = append(reg.Projects, entry)
	if err := storage.WriteProjects(home, reg); err != nil {
		return projectFail(err)
	}
	fmt.Printf("registered %s\n", name)
	return nil
}

func runProjectUnregister(home string, args []string) error {
	remaining, err := lessflags.
		Help("-h,--help", projectUnregisterHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 1 {
		return projectFail(fmt.Errorf("tsk project unregister: name required"))
	}
	name := strings.TrimSpace(remaining[0])
	if name == "" {
		return projectFail(fmt.Errorf("tsk project unregister: name required"))
	}
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return projectFail(err)
	}
	out := reg.Projects[:0]
	found := false
	for _, e := range reg.Projects {
		if e.Name == name {
			found = true
			continue
		}
		out = append(out, e)
	}
	if !found {
		return projectFail(fmt.Errorf("tsk project unregister: unknown project %q", name))
	}
	reg.Projects = out
	if err := storage.WriteProjects(home, reg); err != nil {
		return projectFail(err)
	}
	fmt.Printf("unregistered %s\n", name)
	return nil
}

type projectListJSONRow struct {
	Origin string `json:"origin,omitempty"`
	Name   string `json:"name,omitempty"`
	Cwd    string `json:"cwd,omitempty"`
	Tasks  int    `json:"tasks"`
}

func runProjectRegistryList(home string, args []string) error {
	var asJSON, activeOnly, allFlag, autoOnly, registeredOnly bool
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Bool("--active", &activeOnly).
		Bool("--all", &allFlag).
		Bool("--auto", &autoOnly).
		Bool("--registered", &registeredOnly).
		Help("-h,--help", projectRegistryListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project list: unexpected arguments"))
	}

	modeCount := 0
	if allFlag {
		modeCount++
	}
	if autoOnly {
		modeCount++
	}
	if registeredOnly {
		modeCount++
	}
	if modeCount > 1 {
		return projectFail(fmt.Errorf("tsk project list: --all, --auto, and --registered are mutually exclusive"))
	}

	mode := "all"
	if autoOnly {
		mode = "auto"
	} else if registeredOnly {
		mode = "registered"
	}

	return listProjects(home, mode, asJSON, activeOnly)
}

func listProjects(home, mode string, asJSON, activeOnly bool) error {
	auto, err := storage.ReadProjectsAuto(home)
	if err != nil {
		return projectFail(err)
	}
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return projectFail(err)
	}
	counts, err := countTasksByProject(home)
	if err != nil {
		return projectFail(err)
	}

	includeTasks := mode != "registered" || activeOnly
	var rows []projectListJSONRow

	switch mode {
	case "registered":
		for _, e := range storage.SortedProjects(reg) {
			key := storage.ProjectAutoKey(e.Origin, e.Name)
			n := counts[key]
			if activeOnly && n == 0 {
				continue
			}
			row := projectListJSONRow{Origin: e.Origin, Name: e.Name, Cwd: e.Cwd}
			if includeTasks {
				row.Tasks = n
			}
			rows = append(rows, row)
		}
	case "auto":
		rows = buildAutoListRows(auto, reg, counts, activeOnly)
	default: // all
		rows = buildAllListRows(auto, reg, counts, activeOnly)
	}

	if includeTasks {
		sortProjectListByTasksDesc(rows)
	}

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		if mode == "registered" && !activeOnly {
			return enc.Encode(storage.ProjectsFile{Projects: storage.SortedProjects(reg)})
		}
		return enc.Encode(map[string]any{"projects": rows})
	}
	if len(rows) == 0 {
		fmt.Println("0 projects")
		return nil
	}
	fmt.Print(formatProjectListTable(rows, includeTasks))
	printProjectCount(len(rows))
	return nil
}

func buildAutoListRows(auto storage.ProjectsAutoFile, reg storage.ProjectsFile, counts map[string]int, activeOnly bool) []projectListJSONRow {
	var rows []projectListJSONRow
	for _, e := range storage.SortedProjectsAuto(auto) {
		name := displayNameForAutoEntry(e, reg)
		key := storage.ProjectAutoKey(e.Origin, e.Name)
		n := counts[key]
		if activeOnly && n == 0 {
			continue
		}
		rows = append(rows, projectListJSONRow{
			Origin: e.Origin,
			Name:   name,
			Cwd:    e.Cwd,
			Tasks:  n,
		})
	}
	return rows
}

func buildAllListRows(auto storage.ProjectsAutoFile, reg storage.ProjectsFile, counts map[string]int, activeOnly bool) []projectListJSONRow {
	seen := make(map[string]bool)
	var rows []projectListJSONRow

	for _, e := range storage.SortedProjectsAuto(auto) {
		key := storage.ProjectAutoKey(e.Origin, e.Name)
		seen[key] = true
		name := displayNameForAutoEntry(e, reg)
		n := counts[key]
		if activeOnly && n == 0 {
			continue
		}
		rows = append(rows, projectListJSONRow{
			Origin: e.Origin,
			Name:   name,
			Cwd:    e.Cwd,
			Tasks:  n,
		})
	}

	for _, e := range storage.SortedProjects(reg) {
		key := storage.ProjectAutoKey(e.Origin, e.Name)
		if key == "" || seen[key] {
			continue
		}
		// Also skip if registry name aliases an origin already listed.
		if e.Origin == "" && e.Name != "" {
			skip := false
			for _, r := range rows {
				if r.Name == e.Name {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
		}
		n := counts[key]
		if activeOnly && n == 0 {
			continue
		}
		rows = append(rows, projectListJSONRow{
			Origin: e.Origin,
			Name:   e.Name,
			Cwd:    e.Cwd,
			Tasks:  n,
		})
		seen[key] = true
	}

	sort.Slice(rows, func(i, j int) bool {
		ni, nj := rows[i].Name, rows[j].Name
		if ni != nj {
			if ni == "" {
				return false
			}
			if nj == "" {
				return true
			}
			return ni < nj
		}
		return rows[i].Origin < rows[j].Origin
	})
	return rows
}

func displayNameForAutoEntry(e storage.ProjectAutoEntry, reg storage.ProjectsFile) string {
	if e.Name != "" {
		return e.Name
	}
	if e.Origin == "" {
		return ""
	}
	if re, ok := storage.FindProjectByOrigin(reg, e.Origin); ok && re.Name != "" {
		return re.Name
	}
	_, base, nerr := NormalizeOriginURL("https://" + e.Origin)
	if nerr == nil && base != "" {
		return base
	}
	return filepath.Base(e.Origin)
}

func printProjectCount(n int) {
	word := "projects"
	if n == 1 {
		word = "project"
	}
	fmt.Printf("%d %s\n", n, word)
}

// sortProjectListByTasksDesc orders by tasks descending, then name, then origin.
func sortProjectListByTasksDesc(rows []projectListJSONRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Tasks != rows[j].Tasks {
			return rows[i].Tasks > rows[j].Tasks
		}
		if rows[i].Name != rows[j].Name {
			return rows[i].Name < rows[j].Name
		}
		return rows[i].Origin < rows[j].Origin
	})
}

func countTasksByProject(home string) (map[string]int, error) {
	entries, err := loadProjectTasks(home)
	if err != nil {
		return nil, err
	}
	out := make(map[string]int)
	for _, e := range entries {
		if e.Task.Project == nil {
			continue
		}
		key := storage.ProjectAutoKey(e.Task.Project.Origin, e.Task.Project.Name)
		if key == "" {
			continue
		}
		out[key]++
	}
	return out, nil
}

type projectTreeJSON struct {
	Projects []projectTreeJSONProject `json:"projects"`
}

type projectTreeJSONProject struct {
	Key   string                  `json:"key"`
	Label string                  `json:"label"`
	Tasks []storage.TopicTreeTask `json:"tasks"`
}

func runProjectTree(home string, args []string) error {
	var nameFlag, projectFlag, stageFlag, dirFlag string
	var allFlag, doneFlag, archivedFlag, asJSON, colorFlag, plain bool
	remaining, err := lessflags.
		String("--name", &nameFlag).
		String("--project", &projectFlag).
		String("--dir", &dirFlag).
		String("--stage", &stageFlag).
		Bool("--all", &allFlag).
		Bool("--done", &doneFlag).
		Bool("--archived", &archivedFlag).
		Bool("--json", &asJSON).
		Bool("--color", &colorFlag).
		Bool("--plain", &plain).
		Help("-h,--help", projectTreeHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project tree: unexpected arguments"))
	}
	if nameFlag != "" && projectFlag != "" {
		return projectFail(fmt.Errorf("tsk project tree: --name conflicts with --project"))
	}
	if stageFlag != "" && (doneFlag || archivedFlag) {
		return projectFail(fmt.Errorf("tsk project tree: --stage conflicts with --done/--archived"))
	}
	if strings.TrimSpace(dirFlag) != "" && (nameFlag != "" || projectFlag != "" || allFlag) {
		return projectFail(fmt.Errorf("tsk project tree: --dir conflicts with --name/--project/--all"))
	}

	reg, err := storage.ReadProjects(home)
	if err != nil {
		return projectFail(err)
	}

	wantKey := strings.TrimSpace(projectFlag)
	wantName := strings.TrimSpace(nameFlag)
	var emptyBranch *projectGroup
	if wantKey == "" && wantName == "" && !allFlag {
		probeDir, err := resolveProbeDir(dirFlag)
		if err != nil {
			return projectFail(err)
		}
		origin, _ := gitNormalizedOrigin(probeDir)
		if origin != "" {
			wantKey = origin
			label := displayLabelForOrigin(reg, origin)
			emptyBranch = &projectGroup{Key: origin, Label: label}
		} else if entry, ok := findRegistryByAbsCwd(reg, probeDir); ok && entry.Name != "" {
			wantName = entry.Name
			emptyBranch = &projectGroup{Key: "name:" + entry.Name, Label: entry.Name}
		} else {
			return projectFail(fmt.Errorf("tsk project tree: no git origin and not registered (see tsk project register --help)"))
		}
	}

	entries, err := loadProjectTasks(home)
	if err != nil {
		return err
	}

	var filtered []projectTaskEntry
	for _, e := range entries {
		if !taskMatchesProjectFilter(e.Task, wantKey, wantName, reg) {
			continue
		}
		if !projectTreeStageAllowed(e.Task.Stage, stageFlag, doneFlag, archivedFlag, allFlag) {
			continue
		}
		filtered = append(filtered, e)
	}

	groups := groupProjectTasks(filtered, reg)
	if wantName != "" && wantKey == "" {
		keys := make([]string, 0, len(groups))
		for k := range groups {
			keys = append(keys, k)
		}
		if len(keys) > 1 {
			fmt.Fprintf(os.Stderr, "warning: name %q matches %d project keys; pass --project to disambiguate\n", wantName, len(keys))
		}
	}

	if emptyBranch != nil && len(groups) == 0 {
		groups = map[string]*projectGroup{emptyBranch.Key: emptyBranch}
	}

	color := treeColorEnabled(colorFlag, plain, asJSON)
	if asJSON {
		out := projectTreeJSON{Projects: make([]projectTreeJSONProject, 0, len(groups))}
		for _, g := range sortedProjectGroups(groups) {
			out.Projects = append(out.Projects, projectTreeJSONProject{
				Key:   g.Key,
				Label: g.Label,
				Tasks: buildProjectTaskTree(g.Tasks),
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(out)
	}
	fancy := !plain && isStdoutTTY()
	fmt.Print(formatProjectTree(groups, color, fancy))
	return nil
}

// projectTreeStageAllowed applies project-tree stage visibility:
//   - default: non-terminal only
//   - --done / --archived: only those stages (union when both)
//   - --all alone: all stages
//   - --all with --done/--archived: all projects already selected; stage flags narrow
//   - --stage: exact stage
func projectTreeStageAllowed(stage, stageFlag string, doneFlag, archivedFlag, allFlag bool) bool {
	if stageFlag != "" {
		return stage == stageFlag
	}
	if doneFlag || archivedFlag {
		if doneFlag && stage == "done" {
			return true
		}
		if archivedFlag && stage == "archived" {
			return true
		}
		return false
	}
	if allFlag {
		return true
	}
	return !storage.IsTerminal(stage)
}

func taskMatchesProjectFilter(task storage.Task, wantKey, wantName string, reg storage.ProjectsFile) bool {
	if task.Project == nil {
		return false
	}
	if wantKey != "" {
		if task.Project.Origin != "" {
			return task.Project.Origin == wantKey || filepath.Base(task.Project.Origin) == wantKey
		}
		return task.Project.Name == wantKey
	}
	if wantName != "" {
		if task.Project.Name == wantName {
			return true
		}
		if task.Project.Origin != "" {
			if e, ok := storage.FindProjectByOrigin(reg, task.Project.Origin); ok && e.Name == wantName {
				return true
			}
			_, base, _ := NormalizeOriginURL("https://" + task.Project.Origin)
			return base == wantName || filepath.Base(task.Project.Origin) == wantName
		}
		return false
	}
	return true
}

type projectTaskEntry struct {
	Task storage.Task
	Dir  string
}

type projectGroup struct {
	Key   string
	Label string
	Tasks []projectTaskEntry
}

func loadProjectTasks(home string) ([]projectTaskEntry, error) {
	ids, err := storage.ListTaskIDs(home)
	if err != nil {
		return nil, err
	}
	var out []projectTaskEntry
	for _, id := range ids {
		task, dir, err := storage.LoadTaskByID(home, id)
		if err != nil {
			return nil, err
		}
		if task.Project == nil {
			continue
		}
		if task.Project.Origin == "" && task.Project.Name == "" {
			continue
		}
		out = append(out, projectTaskEntry{Task: task, Dir: dir})
	}
	return out, nil
}

func projectGroupKey(ref *storage.ProjectRef) string {
	if ref.Origin != "" {
		return ref.Origin
	}
	return "name:" + ref.Name
}

func groupProjectTasks(entries []projectTaskEntry, reg storage.ProjectsFile) map[string]*projectGroup {
	groups := make(map[string]*projectGroup)
	for _, e := range entries {
		key := projectGroupKey(e.Task.Project)
		g, ok := groups[key]
		if !ok {
			g = &projectGroup{
				Key:   key,
				Label: displayLabelForRef(reg, e.Task.Project),
			}
			groups[key] = g
		}
		g.Tasks = append(g.Tasks, e)
	}
	return groups
}

func displayLabelForRef(reg storage.ProjectsFile, ref *storage.ProjectRef) string {
	if ref.Origin != "" {
		return displayLabelForOrigin(reg, ref.Origin)
	}
	return ref.Name
}

func displayLabelForOrigin(reg storage.ProjectsFile, origin string) string {
	name := ""
	if e, ok := storage.FindProjectByOrigin(reg, origin); ok && e.Name != "" {
		name = e.Name
	} else {
		_, base, err := NormalizeOriginURL("https://" + origin)
		if err == nil && base != "" {
			name = base
		} else {
			name = filepath.Base(origin)
		}
	}
	return name + "  " + origin
}

func sortedProjectGroups(groups map[string]*projectGroup) []*projectGroup {
	out := make([]*projectGroup, 0, len(groups))
	for _, g := range groups {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Label != out[j].Label {
			return out[i].Label < out[j].Label
		}
		return out[i].Key < out[j].Key
	})
	return out
}

func buildProjectTaskTree(entries []projectTaskEntry) []storage.TopicTreeTask {
	byID := make(map[int]projectTaskEntry, len(entries))
	for _, e := range entries {
		byID[e.Task.ID] = e
	}
	children := make(map[int][]int)
	var roots []int
	for _, e := range entries {
		parent := e.Task.ParentID
		if parent != 0 {
			if _, ok := byID[parent]; ok {
				children[parent] = append(children[parent], e.Task.ID)
				continue
			}
		}
		roots = append(roots, e.Task.ID)
	}
	sort.Ints(roots)
	for id := range children {
		sort.Ints(children[id])
	}
	var build func(id int) storage.TopicTreeTask
	build = func(id int) storage.TopicTreeTask {
		e := byID[id]
		node := storage.TopicTreeTask{
			ID:    e.Task.ID,
			Stage: e.Task.Stage,
			Slug:  e.Task.Slug,
			Title: e.Task.Title,
			Dir:   filepath.Base(e.Dir),
		}
		for _, cid := range children[id] {
			node.Tasks = append(node.Tasks, build(cid))
		}
		return node
	}
	out := make([]storage.TopicTreeTask, 0, len(roots))
	for _, id := range roots {
		out = append(out, build(id))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Dir < out[j].Dir })
	return out
}

func formatProjectTree(groups map[string]*projectGroup, color, fancy bool) string {
	var b strings.Builder
	b.WriteString(".\n")
	ordered := sortedProjectGroups(groups)
	if len(ordered) == 0 {
		b.WriteString("(empty)\n")
		b.WriteString("0 tasks, 0 projects\n")
		return b.String()
	}

	rootKids := make([]*renderNode, 0, len(ordered))
	taskCount := 0
	mark := projectMark(fancy)
	idWidth := 0
	trees := make([][]storage.TopicTreeTask, 0, len(ordered))
	for _, g := range ordered {
		tree := buildProjectTaskTree(g.Tasks)
		trees = append(trees, tree)
		if w := maxTaskIDWidth(tree); w > idWidth {
			idWidth = w
		}
	}
	for i, g := range ordered {
		tree := trees[i]
		taskCount += countProjectTreeTasks(tree)
		short, origin := splitProjectDisplayLabel(g.Label)
		node := &renderNode{
			name:  mark + short,
			extra: formatProjectOriginExtra(origin, color),
		}
		node.children = projectTasksToRender(tree, color, idWidth)
		rootKids = append(rootKids, node)
	}
	writeRenderKids(&b, rootKids, "")

	taskWord := "tasks"
	if taskCount == 1 {
		taskWord = "task"
	}
	projectWord := "projects"
	if len(ordered) == 1 {
		projectWord = "project"
	}
	b.WriteString(fmt.Sprintf("\n%d %s, %d %s\n", taskCount, taskWord, len(ordered), projectWord))
	return b.String()
}

func countProjectTreeTasks(tasks []storage.TopicTreeTask) int {
	n := 0
	for i := range tasks {
		n += countTaskNode(&tasks[i])
	}
	return n
}

func projectTasksToRender(tasks []storage.TopicTreeTask, color bool, idWidth int) []*renderNode {
	out := make([]*renderNode, 0, len(tasks))
	for _, task := range tasks {
		node := &renderNode{
			name:  formatTaskLeafName(task.ID, task.Title, task.Slug, idWidth),
			extra: formatTaskStageExtra(task.Stage, color),
			style: taskStageStyle(task.Stage, color),
		}
		if len(task.Tasks) > 0 {
			node.children = projectTasksToRender(task.Tasks, color, idWidth)
		}
		out = append(out, node)
	}
	return out
}

func resolveProbeDir(dirFlag string) (string, error) {
	if strings.TrimSpace(dirFlag) == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getwd: %w", err)
		}
		return filepath.Abs(wd)
	}
	abs, err := filepath.Abs(dirFlag)
	if err != nil {
		return "", fmt.Errorf("resolve --dir: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("resolve --dir: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("--dir is not a directory: %s", abs)
	}
	return abs, nil
}

func gitNormalizedOrigin(dir string) (string, error) {
	raw, err := gitRemoteOriginURL(dir)
	if err != nil || strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("no origin")
	}
	id, _, err := NormalizeOriginURL(raw)
	if err != nil {
		return "", err
	}
	return id, nil
}

func gitRemoteOriginURL(dir string) (string, error) {
	top := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	if out, err := top.Output(); err != nil || strings.TrimSpace(string(out)) == "" {
		return "", fmt.Errorf("not a git repo")
	}
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func findRegistryByAbsCwd(reg storage.ProjectsFile, abs string) (storage.ProjectEntry, bool) {
	abs = filepath.Clean(abs)
	for _, e := range reg.Projects {
		if e.Cwd == "" {
			continue
		}
		expanded := filepath.Clean(pathfmt.Expand(e.Cwd))
		if expanded == abs {
			return e, true
		}
	}
	return storage.ProjectEntry{}, false
}

// mainRepoDirForAuto returns the main worktree root for git dirs, else dir.
func mainRepoDirForAuto(dir string) string {
	if main, err := gitMainRepoDir(dir); err == nil && main != "" {
		return main
	}
	return dir
}

// gitMainRepoDir resolves the primary worktree directory (not a linked worktree).
func gitMainRepoDir(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-common-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	common := strings.TrimSpace(string(out))
	if common == "" {
		return "", fmt.Errorf("empty git-common-dir")
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(dir, common)
	}
	common, err = filepath.Abs(common)
	if err != nil {
		return "", err
	}
	common = filepath.Clean(common)
	if filepath.Base(common) == ".git" {
		return filepath.Dir(common), nil
	}
	// Bare or unusual layout: fall back to show-toplevel.
	top := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	tout, terr := top.Output()
	if terr != nil {
		return "", terr
	}
	return filepath.Clean(strings.TrimSpace(string(tout))), nil
}

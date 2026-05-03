# Implementation Plan: Course Hooks

## Objective
Implement a hook system that allows a custom function to run on a course folder just before it gets processed. The hooks are identified by the course's ID.

## Key Files & Context
- `internal/hooks/hooks.go` (New File): Contains the hook registry and execution logic.
- `internal/service/service.go`: Update `Process` and `ExtractAndProcess` to accept a course ID and run the hook.
- `internal/service/course.go`: Pass the course ID down to the processing functions.

## Implementation Steps
1. **Create `internal/hooks/hooks.go`**
   - Define `type HookFunc func(folderPath string) error`.
   - Initialize `var Registry = map[string]HookFunc{}`.
   - Add `func Run(courseID string, folderPath string) error`:
     - Checks if `Registry[courseID]` exists.
     - If it exists, calls it with `folderPath` and returns the error (if any).
     - If it doesn't exist, returns `nil`.

2. **Update `internal/service/service.go`**
   - Modify `func Process(courseID string, source string, targets []string, cfg *config.Config) error`.
   - At the beginning of `Process`, add:
     ```go
     if err := hooks.Run(courseID, source); err != nil {
         return err
     }
     ```
   - Import `github.com/dhanush-cache/course-flow/internal/hooks`.
   - Modify `func ExtractAndProcess(courseID string, zipFiles []string, targets []string, cfg *config.Config) error` to accept `courseID` and pass it to `Process`.

3. **Update `internal/service/course.go`**
   - In `func AddCourse(...)`, update the calls to `ExtractAndProcess` and `Process`:
     - `err = ExtractAndProcess(course.ID, options.ZipFiles, fileNames, cfg)`
     - `err = Process(course.ID, options.Directory, fileNames, cfg)`

## Verification & Testing
- Build the project using `go build .` to ensure there are no compilation errors.
- Verify that standard processing paths without hooks continue to work seamlessly.
- (Optional/Manual) Add a mock hook temporarily to test execution on a specific test slug.

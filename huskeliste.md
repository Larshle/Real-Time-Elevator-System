

- [ ] Direction på oppstart av heis
- [ ] dobbelsjekke packetloss
- [ ] teste mer på isolated
- [ ] watchdog
- [ ] readme
- [ ] fjerne unødvendige filer (f.eks packetloss)
- [ ] rydde i kode, spesielt assigner
- [ ] fjerne printing
- [ ] ikke ha med EmptyAssigner
- [ ] huske å endre doopOpen til 3 sekunder


### Main Requirements Checklist

#### Button Lights as Service Guarantee
- [X] Implement hall call button functionality to ensure an elevator arrives once activated.
- [X] Implement cab call button functionality, with the call being specific to the elevator in the workspace.

#### No Calls Lost
- [ ] Ensure the system handles failure states without losing calls.
- [ ] Implement measures for cab calls to be executed once service is restored, even after power or software failures.
- [X] Allow for network-disconnected elevators to serve current and new cab calls without needing reinitialization.

#### Lights and Buttons Functionality
- [X] Enable hall call buttons to summon elevators in all workspaces.
- [X] Ensure hall button lights reflect the same state across workspaces under normal circumstances.
- [X] Keep cab button lights separate for each workspace.
- [X] Ensure button lights activate reasonably quickly after being pressed and deactivate once the call is serviced.

#### Door Functionality
- [X] Use the "door open" lamp to indicate door status.
- [X] Ensure the door (light) is not on while the elevator is moving.
- [X] Set door open duration to 3 seconds when stopping at a floor.
- [X] Implement obstruction switch functionality to prevent door closing when obstructed.

#### Individual Elevator Efficiency
- [X] Avoid unnecessary stops at every floor.
- [X] Manage clearing of hall call buttons correctly, indicating the elevator's direction and not clearing up and down calls simultaneously.
- [X] Allow the elevator to change announced direction efficiently if needed.

### Secondary Requirements Checklist

#### Efficiency in Call Serving
- [X] Distribute calls across elevators to ensure they are serviced promptly.

### Permitted Assumptions
- [X] At least one elevator is always operational.
- [X] Single elevator or disconnected elevator cab call redundancy is not required.
- [X] No network partitioning will occur.

### Unspecified Behavior
- [ ] Decide on behavior during network disconnection at initialization.
- [ ] Optional implementation of new hall calls when the elevator is network-disconnected.
- [ ] Define the functionality of the stop button, if implemented.

### Recommendations
- [ ] Start development with 1 to 3 elevators and 4 floors, avoiding hard-coded values.
- [ ] Implement a command-line switch for specifying elevator identifiers (--id <number>).

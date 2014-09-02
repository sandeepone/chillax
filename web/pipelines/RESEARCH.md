* A pipeline consists of multiple stages.

* Each pipeline is handled by 1 goroutine.

* Each pipeline definition is serialized to storage.

* Only 1 chillax host can handle the whole lifecycle of a pipeline.

* If a stage is fanning out to N stages, that stage creates a pool of N goroutines.

* If a stage is failed to complete, its pipeline will try again on 2^^n interval.

* When all possible leaf stages failed, stop the pipeline and save its state on storage.

* Each stage handles its own params. Either as form data or json body. It really depends on the underlying HTTP endpoint.

* Each stage has errChan and outChan.

* Each stage failure is serialized on storage.

* Each stage success is recorded as LastSuccessAt on storage.

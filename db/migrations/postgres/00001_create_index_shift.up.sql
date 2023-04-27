CREATE INDEX shift_start_idx ON public."Shift" ("start", "end");
CREATE INDEX shift_essential_idx ON public."Shift" ("worker_id", "start", "end", "profession", "facility_id");
CREATE INDEX facility_is_active_idx ON public."Facility" (is_active);
CREATE INDEX facilityrequirement_facility_id_idx ON public."FacilityRequirement" (facility_id);
CREATE INDEX documentworker_worker_id_idx ON public."DocumentWorker" (worker_id);
ALTER TABLE module_info ADD CONSTRAINT module_info_updated CHECK (updated_at >= created_at);
ALTER TABLE module_info ADD CONSTRAINT module_info_duration CHECK (moduleDuration BETWEEN 5 AND 15);


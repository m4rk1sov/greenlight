CREATE TABLE module_info (
                             id BIGSERIAL PRIMARY KEY,
                             created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                             moduleName VARCHAR(255) NOT NULL,
                             moduleDuration INTEGER NOT NULL,
                             examType VARCHAR(255) NOT NULL,
                             version INTEGER NOT NULL DEFAULT 1
);
\c proyecto_recolecta;

-- =====================
-- CONSTRAINTS
-- =====================

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_rol_nombre'
    ) THEN
        ALTER TABLE rol ADD CONSTRAINT uq_rol_nombre UNIQUE (nombre);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_rol_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT fk_rol FOREIGN KEY (rol_id) REFERENCES rol(id);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_mail_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT uq_mail_empleado UNIQUE (mail);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_username_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT uq_username_empleado UNIQUE (username);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_empleado_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT fk_empleado_licencia FOREIGN KEY (id_empleado) REFERENCES empleado(id);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT uq_licencia UNIQUE (licencia);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_vencimiento'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT chk_fecha_vencimiento CHECK (fecha_vencimiento > CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_empleado_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT fk_empleado_licencia FOREIGN KEY (id_empleado) REFERENCES empleado(id);
    END IF;
END $$;


DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_empleado_historial'
    ) THEN
        ALTER TABLE historial_asignacion ADD CONSTRAINT fk_empleado_historial FOREIGN KEY (id_empleado) REFERENCES empleado(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_historial'
    ) THEN
        ALTER TABLE historial_asignacion ADD CONSTRAINT fk_camion_historial FOREIGN KEY (id_camion) REFERENCES camion(camion_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_empleado_camion_historial'
    ) THEN
        ALTER TABLE historial_asignacion ADD CONSTRAINT uq_empleado_camion_historial UNIQUE (id_empleado, id_camion, fecha_asignacion);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_asignacion'
    ) THEN
        ALTER TABLE historial_asignacion ADD CONSTRAINT chk_fecha_asignacion CHECK (fecha_asignacion <= CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_baja'
    ) THEN
        ALTER TABLE historial_asignacion ADD CONSTRAINT chk_fecha_baja CHECK (fecha_baja IS NULL OR fecha_baja >= fecha_asignacion);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_placa_camion'
    ) THEN
        ALTER TABLE camion ADD CONSTRAINT uq_placa_camion UNIQUE (placa);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_tipo_camion'
    ) THEN
        ALTER TABLE camion ADD CONSTRAINT fk_tipo_camion FOREIGN KEY (tipo_camion_id) REFERENCES tipo_camion(tipo_camion_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_estado_camion_camion'
    ) THEN
        ALTER TABLE estado_camion ADD CONSTRAINT fk_estado_camion_camion FOREIGN KEY (camion_id) REFERENCES camion(camion_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_afectado'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT fk_camion_afectado FOREIGN KEY (camion_id) REFERENCES camion(camion_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_tipo_mantenimiento'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT fk_tipo_mantenimiento FOREIGN KEY (tipo_mantenimiento) REFERENCES tipo_mantenimiento(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_kilometraje_mantenimiento'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT chk_kilometraje_mantenimiento CHECK (kilometraje >= 0);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_asignado_ruta'
    ) THEN
        ALTER TABLE registro_asignacion_ruta ADD CONSTRAINT fk_camion_asignado_ruta FOREIGN KEY (camion_id) REFERENCES camion(camion_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_asignacion_ruta'
    ) THEN
        ALTER TABLE registro_asignacion_ruta ADD CONSTRAINT chk_fecha_asignacion_ruta CHECK (fecha_asignacion <= CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_asignada'
    ) THEN
        ALTER TABLE registro_asignacion_ruta ADD CONSTRAINT fk_ruta_asignada FOREIGN KEY (ruta_id) REFERENCES ruta(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_colonia_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT fk_colonia_ruta FOREIGN KEY (colonia_id) REFERENCES colonia(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_nombre_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT chk_nombre_ruta CHECK (nombre <> '');
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'uq_nombre_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT uq_nombre_ruta UNIQUE (nombre);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_id_punto_recoleccion'
    ) THEN
        ALTER TABLE punto_recoleccion ADD CONSTRAINT fk_ruta_id_punto_recoleccion FOREIGN KEY (ruta_id) REFERENCES ruta(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_email_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT uq_email_ciudadano UNIQUE (email);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_alias_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT uq_alias_ciudadano UNIQUE (alias);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_email_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT chk_email_ciudadano CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ciudadano_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT fk_ciudadano_domicilio FOREIGN KEY (ciudadano_id) REFERENCES ciudadano(id);
    END IF;

    if NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_colonia_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT fk_colonia_domicilio FOREIGN KEY (colonia_id) REFERENCES colonia(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_alias_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT chk_alias_domicilio CHECK (alias <> '');
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'domicilio' AND column_name = 'calle'
    ) THEN
        ALTER TABLE domicilio ADD COLUMN calle VARCHAR(255);
        ALTER TABLE domicilio ADD COLUMN numero VARCHAR(50) DEFAULT '';
        ALTER TABLE domicilio ADD COLUMN referencia VARCHAR(255);
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'domicilio' AND column_name = 'direccion'
        ) THEN
            UPDATE domicilio SET calle = direccion WHERE calle IS NULL;
        END IF;
        UPDATE domicilio SET calle = '' WHERE calle IS NULL;
        UPDATE domicilio SET numero = '' WHERE numero IS NULL;
        ALTER TABLE domicilio ALTER COLUMN calle SET NOT NULL;
        ALTER TABLE domicilio ALTER COLUMN numero SET NOT NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'domicilio' AND column_name = 'latitud'
    ) THEN
        ALTER TABLE domicilio ADD COLUMN latitud DOUBLE PRECISION;
        ALTER TABLE domicilio ADD COLUMN longitud DOUBLE PRECISION;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'domicilio' AND column_name = 'direccion'
          AND is_generated = 'NEVER'
    ) THEN
        UPDATE domicilio
        SET calle = COALESCE(NULLIF(TRIM(calle), ''), direccion)
        WHERE calle IS NULL OR TRIM(calle) = '';
        ALTER TABLE domicilio DROP COLUMN direccion;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'domicilio' AND column_name = 'direccion'
    ) THEN
        ALTER TABLE domicilio ADD COLUMN direccion VARCHAR(255)
            GENERATED ALWAYS AS (TRIM(BOTH FROM calle || ' ' || COALESCE(numero, ''))) STORED;
    END IF;

END $$;

package oracle

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone_status (
            entityone_status_id NUMBER(10,0) NOT NULL,
            action_id NUMBER(3, 0) NOT NULL,
            status_id NUMBER(3, 0) NOT NULL ,
            time_created DATE DEFAULT SYSDATE NOT NULL,
			CONSTRAINT es_pk PRIMARY KEY (entityone_status_id)
				USING INDEX (CREATE INDEX es_pk_esi ON entityone_status(entityone_status_id))			
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone_status %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE SEQUENCE entityone_status_seq START WITH 1`,
	)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		CREATE OR REPLACE TRIGGER entityone_status_trig
		BEFORE INSERT ON entityone_status FOR EACH ROW
		BEGIN
			SELECT entityone_status_seq.NEXTVAL
			INTO   :new.entityone_status_id
			FROM   dual;
		END;
	`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create trigger %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_si ON entityone_status(status_id)`,
	)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create trigger %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone (
            entityone_id NUMBER(10,0) NOT NULL,
			time_created DATE DEFAULT SYSDATE NOT NULL,
			entityone_status_id  NUMBER(10,0) NOT NULL,
			CONSTRAINT e_pk PRIMARY KEY (entityone_id)
				USING INDEX (CREATE INDEX e_pk_ei ON entityone(entityone_id)),
            CONSTRAINT es_fk_esi FOREIGN KEY (entityone_status_id)
            	REFERENCES entityone_status(entityone_status_id)
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE SEQUENCE entityone_seq START WITH 1`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		CREATE OR REPLACE TRIGGER entityone_trig
		BEFORE INSERT ON entityone FOR EACH ROW
		BEGIN
			SELECT entityone_seq.NEXTVAL
			INTO   :new.entityone_id
			FROM   dual;
		END;
	`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE UNIQUE INDEX e_idx_esi ON entityone(entityone_status_id)`,
	)

	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'ENTITYONE';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP TABLE ENTITYONE';
			END IF;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'ENTITYONE_STATUS';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP TABLE ENTITYONE_STATUS';
			END IF;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_sequences WHERE sequence_name = 'ENTITYONE_STATUS_SEQ';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP SEQUENCE entityone_status_seq';
			END IF;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_sequences WHERE sequence_name = 'ENTITYONE_SEQ';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP SEQUENCE entityone_seq';
			END IF;
		END;
	`)

	return errExec
}

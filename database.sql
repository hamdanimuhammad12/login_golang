CREATE TABLE public.users (
	id int4 NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
	"name" varchar(255) NULL,
	phone varchar(15) NULL,
	username varchar(255) NULL,
	"password" varchar(255) NULL,
	role_id int2 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	is_active bool NULL,
	deleted_at timestamp NULL,
	email varchar NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
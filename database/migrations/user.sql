CREATE TABLE public.users (
	uuid uuid NOT NULL,
	created_at date NOT NULL,
	updated_at date NULL,
	email varchar(50) NOT NULL,
	username varchar(20) NOT NULL,
	"password" varchar(255) NOT NULL,
	deleted_at date NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (uuid),
	CONSTRAINT users_username_key UNIQUE (username)
);
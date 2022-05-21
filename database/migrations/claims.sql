CREATE TABLE public.claims (
	audience text NULL,
	expires_at int8 NULL,
	id bigserial NOT NULL,
	issued_at int8 NULL,
	issuer text NULL,
	not_before int8 NULL,
	subject text NULL,
	CONSTRAINT claims_pkey PRIMARY KEY (id)
);
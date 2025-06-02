-- name: GetAllLanguages :many
select
  *
from
  languages;

-- name: GetSubmissionByCredentials :one
select
  *
from
  submissions
where
  username = $1
  and password = $2;

-- name: CreateSubmission :one
insert 
  into submissions 
    (name, phone, email, birth_date, bio, sex, username, password)
  values 
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetSubmissionByIdAndCredentials :one
select
  id
from
  submissions
where
  id = $1
  and username = $2
  and password = $3;

-- name: GetFullSubmissionByIdAndCredentials :one
select * from submissions where id = $1 and username = $2 and password = $3;

-- name: GetSubmissionLanguages :many
select
 language_id 
from
  submission_languages
where
  submission_id = $1;

-- name: UpdateSubmission :exec
update
  submissions
set
  name = $1,
  phone = $2,
  email = $3,
  birth_date = $4,
  bio = $5,
  sex = $6
where
  id = $7
returning *;

-- name: AddLanguageToSubmission :exec
insert 
  into submission_languages 
    (submission_id, language_id)
  values 
    ($1, $2);

-- name: DeleteSubmissionLanguages :exec
delete from submission_languages where submission_id = $1;
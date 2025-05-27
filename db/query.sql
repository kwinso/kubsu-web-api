-- name: GetAllLanguages :many
select
  id
from
  languages;


-- name: CreateSubmission :one
insert 
  into submissions 
    (name, phone, email, birth_date, bio, sex, username, password)
  values 
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: AddLanguageToSubmission :exec
insert 
  into submission_languages 
    (submission_id, language_id)
  values 
    ($1, $2);
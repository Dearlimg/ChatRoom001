-- name: CreateApplication :exec
insert into applications (applications.account1_id,applications.account2_id,applications.apply_msg,applications.refuse_msg)
values ( ?,? ,?,'');

-- name: CreateGet :one
select * from applications where account1_id = ?;

-- name: ExistsApplicationByIDWithLock :one
select exists(
    select 1
    from applications
    where (account2_id = ? and account1_id = ?)
       or (account1_id = ? and account2_id = ?)
        for update
);

-- name: DeleteApplication :exec
delete
from applications
where account1_id = ?
  and account2_id = ?;

-- name: GetApplicationByID :one
select *
from applications
where account1_id = ?
  and account2_id = ?
limit  1;

-- name: UpdateApplication :exec
update applications
set status = ?,
    refuse_msg = ?
where account1_id =?
  and account2_id = ?;

-- name: GetApplications :many
select app.*,
       a1.name as account1_name,
       a1.avatar as account1_avatar,
       a2.name as account2_name,
       a2.avatar as account2_avatar
from accounts a1,
     accounts a2,
     (select *, count(*) over () as total
      from applications
      where account1_id=?
         or account2_id = ?
      order by create_at desc
      limit ? offset ?) as app
where a1.id = app.account1_id
  and a2.id = app.account2_id;



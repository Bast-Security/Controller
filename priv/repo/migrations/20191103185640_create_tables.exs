defmodule Bast.Controller.Repo.Migrations.CreateTables do
  use Ecto.Migration

  def change do
    create table(:doors) do
      add :name, :string, size: 32, null: false
    end
    unique_index(:doors, [:name])

    create table(:users) do
      add :name,    :string, size: 32, null: false
      add :email,   :string,           null: false
      add :pin,     :string
      add :cardno,  :integer
      add :lastaccess, references(:doors)
    end
    unique_index(:users, [:email, :phoneno, :cardno])

    create table(:roles) do
      add :name, :string, null: false
    end

    create table(:userrole) do
      add :user, references(:roles), null: false
      add :role, references(:users), null: false
    end

    create table(:activetimes) do
      add :start, :utc_datetime, null: false
      add :end,   :utc_datetime, null: false
    end

    create table(:permissions) do
      add :role, references(:roles), null: false
    end

    create table(:credentialtypes) do
      add :permission, references(:permissions), null: false
      add :type,       :integer,                 null: false
    end

    create table(:permissiontarget) do
      add :door, references(:doors),             null: false
      add :permission, references(:permissions), null: false
    end
  end
end

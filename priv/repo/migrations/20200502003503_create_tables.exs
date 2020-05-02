defmodule Bast.Repo.Migrations.CreateTables do
  use Ecto.Migration

  def change do
    create table(:systems) do
      add :name, :string, null: false
      add :totpkey, :binary
    end

    create table(:admins) do
      add :keyx, :binary
      add :keyy, :binary
      add :challenge, :binary
    end

    create table(:admins_systems) do
      add :system_id, references(:systems)
      add :admin_id, references(:admins)
    end

    create table(:users) do
      add :system_id, references(:systems)
      add :name, :string, null: false
      add :email, :string
      add :phone, :string
      add :card, :string
      add :pin, :string
    end

    create table(:roles) do
      add :system_id, references(:systems)
      add :name, :string, null: false
    end

    create table(:users_roles) do
      add :user_id, references(:users)
      add :role_id, references(:roles)
    end

    create table(:doors) do
      add :system_id, references(:systems)
      add :name, :string, null: false
      add :keyx, :binary
      add :keyy, :binary
      add :challenge, :binary
      add :mode, :int
    end

    create table(:roles_doors) do
      add :role_id, references(:roles)
      add :door_id, references(:doors)
    end

    create table(:history) do
      add :door_id, references(:doors)
      add :card, :string
      add :pin, :string
      add :time, :utc_datetime
    end
  end
end

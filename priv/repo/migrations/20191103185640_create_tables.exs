defmodule Bast.Controller.Repo.Migrations.CreateTables do
  use Ecto.Migration

  def change do
    create table(:doors) do
      add :name, :string, size: 32, null: false
    end
    unique_index(:doors, [:name])

    create table(:users) do
      add :name,    :string, size: 32, null: false
      add :email,   :string
      add :phoneno, :string, size: 10
      add :pubkey,  :binary
      add :pin,     :string
      add :cardno,  :integer
      add :lastaccess, references(:doors)
    end
    unique_index(:users, [:email, :phoneno, :cardno])
  end
end

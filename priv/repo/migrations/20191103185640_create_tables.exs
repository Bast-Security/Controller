defmodule Bast.Middleman.Repo.Migrations.CreateTables do
  use Ecto.Migration

  def change do
    create table(:doors) do
      add :name, :string, size: 32, null: false
    end

    create table(:users) do
      add :name,   :string, size: 32, null: false
      add :pubkey, :binary
      add :pin,    :string
      add :cardno, :integer
      add :lastaccess, references(:doors)
    end
  end
end

defmodule Bast.System do
  use Ecto.Schema

  schema "systems" do
    field :name, :string
    field :totpkey, :binary
    many_to_many :admins, Bast.Admin, join_through: "admins_systems"
  end
end

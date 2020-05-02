defmodule Bast.Admin do
  use Ecto.Schema

  schema "admins" do
    field :keyx, :binary
    field :keyy, :binary
    field :challenge, :binary
    many_to_many :systems, Bast.System, join_through: "admins_systems"
  end
end

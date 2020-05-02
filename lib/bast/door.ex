defmodule Bast.Door do
  use Ecto.Schema

  schema "doors" do
    field :name, :string
    field :keyX, :binary
    field :keyY, :binary
    field :challenge, :binary
    field :mode, :integer
    many_to_many :roles, Bast.Role, join_through: "roles_doors"
    has_many :history, Bast.History
  end
end

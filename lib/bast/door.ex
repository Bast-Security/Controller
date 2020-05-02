defmodule Bast.Door do
  use Ecto.Schema

  schema "doors" do
    field :name, :string
    field :keyx, :binary
    field :keyy, :binary
    field :challenge, :binary
    field :mode, :integer
    many_to_many :roles, Bast.Role, join_through: "roles_doors"
    has_many :history, Bast.History
  end
end

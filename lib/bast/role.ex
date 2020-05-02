defmodule Bast.Role do
  use Ecto.Schema

  schema "roles" do
    field :name, :string
    many_to_many :users, Bast.User, join_through: "users_roles"
    many_to_many :doors, Bast.Door, join_through: "roles_doors"
  end
end

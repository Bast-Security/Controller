defmodule Bast.User do
  use Ecto.Schema

  schema "user" do
    field :name, :string
    field :email, :string
    field :phone, :string
    field :card, :string
    field :pin, :string
    many_to_many :roles, Bast.Role, join_through: "users_roles"
  end
end

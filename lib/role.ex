defmodule Bast.Controller.Role do
  use Ecto.Schema
  import Ecto.Changeset

  def validate(role) do
    role
    |> change
    |> validate_required([:name])
  end

  schema "roles" do
    field :name, :string
    many_to_many :users, Bast.Controller.User, join_through: "userrole"
  end
end

defmodule Bast.Controller.User do
  use Ecto.Schema
  import Ecto.Changeset

  @spec validate({map, any} | %{__struct__: atom | %{__changeset__: any}}) :: Ecto.Changeset.t()
  def validate(user) do
    user
    |> change
    |> validate_required([:name, :email])
  end

  schema "users" do
    field   :name,    :string
    field   :email,   :string
    field   :pin,     :string
    field   :cardno,  :integer
    has_one :doors,   Bast.Controller.Door
    many_to_many :roles, Bast.Controller.Role, join_through: "userrole"
  end
end

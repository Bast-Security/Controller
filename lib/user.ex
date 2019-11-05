defmodule Bast.Controller.User do
  use Ecto.Schema
  import Ecto.Changeset

  def validate(user) do
    user
    |> change
    |> validate_required([:name])
  end

  schema "users" do
    field   :name,    :string
    field   :email,   :string
    field   :phoneno, :string
    field   :pubkey,  :binary
    field   :pin,     :string
    field   :cardno,  :integer
    has_one :door,    Bast.Controller.Door
  end
end

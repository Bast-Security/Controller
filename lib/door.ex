defmodule Bast.Controller.Door do
  use Ecto.Schema

  schema "doors" do
    field :name, :string
    belongs_to :user, Bast.Controller.User
  end
end

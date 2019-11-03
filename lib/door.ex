defmodule Bast.Middleman.Door do
  use Ecto.Schema

  schema "doors" do
    field :name, :string
    belongs_to :user, Bast.Middleman.User
  end
end

defmodule Bast.Middleman.User do
  use Ecto.Schema

  schema "users" do
    field :name, :string
    field :pubkey, :binary
    field :pin, :string
    field :cardno, :integer
    has_one :door, Bast.Middleman.Door
  end
end

defmodule Bast.History do
  use Ecto.Schema

  schema "history" do
    field :card, :string
    field :pin, :string
    field :time, :utc_datetime
    belongs_to :door, Bast.Door
  end
end

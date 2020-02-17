defmodule Bast.Controller.ActiveTime do
  use Ecto.Schema
  import Ecto.Changeset

  def validate(activetime) do
    activetime
    |> change
    |> validate_required([:start, :end])
  end

  schema "activetimes" do
    field :start, :utc_datetime
    field :end,   :utc_datetime
  end
end

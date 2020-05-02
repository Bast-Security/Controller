defmodule Bast.Admin do
  use Ecto.Schema
  alias Ecto.Changeset

  schema "admins" do
    field :keyX, :binary
    field :keyY, :binary
    field :challenge, :binary
    many_to_many :systems, Bast.System, join_through: "admins_systems"
  end

  def create(map) do
    Changeset.cast(%Bast.Admin{}, map, [:keyX, :keyY])
    |> Changeset.validate_required([:keyX, :keyY])
  end
end

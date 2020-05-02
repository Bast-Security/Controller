defmodule Bast.Repo do
  use Ecto.Repo,
    otp_app: :bast,
    adapter: Ecto.Adapters.Postgres
end

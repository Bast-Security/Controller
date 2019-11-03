defmodule Bast.Middleman.Repo do
  use Ecto.Repo,
    otp_app: :middleman,
    adapter: Ecto.Adapters.MyXQL
end

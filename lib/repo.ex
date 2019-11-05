defmodule Bast.Controller.Repo do
  use Ecto.Repo,
    otp_app: :controller,
    adapter: Ecto.Adapters.MyXQL
end

terraform {
  required_providers {
    litellm = {
      source = "jimmyflatting/litellm"
    }
  }
}

provider "litellm" {
  api_key = var.litellm_api_key # API key for LiteLLM authentication
}

resource "litellm_model" "gpt4" {
  name           = "gpt-4-custom"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  api_key        = var.openai_api_key

  metadata = {
    description = "Custom GPT-4 model configuration"
    team        = "ml-team"
  }
}

resource "litellm_key" "ml_team_key" {
  key_alias  = "ml-team-key"
  team_id    = "ml-team-123"
  models     = ["gpt-4-custom"]
  max_budget = 100.0

  depends_on = [litellm_model.gpt4]
}

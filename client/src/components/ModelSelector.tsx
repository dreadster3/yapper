import React from "react";
import { ChevronDown, Zap, Brain, Server } from "lucide-react";
import { AppSettings, ModelProvider, Model } from "../types";

interface ModelSelectorProps {
  settings: AppSettings;
  onModelChange: (provider: string, model: string) => void;
}

const getProviderIcon = (providerId: string) => {
  switch (providerId) {
    case "openai":
      return <Brain className="w-4 h-4" />;
    case "anthropic":
      return <Zap className="w-4 h-4" />;
    case "ollama":
      return <Server className="w-4 h-4" />;
    default:
      return <Brain className="w-4 h-4" />;
  }
};

export const ModelSelector: React.FC<ModelSelectorProps> = ({
  settings,
  onModelChange,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);

  const currentProvider = settings.providers.find(
    (p) => p.id === settings.selectedProvider,
  );
  const currentModel = currentProvider?.models.find(
    (m) => m.id === settings.selectedModel,
  );

  const handleModelSelect = (provider: ModelProvider, model: Model) => {
    onModelChange(provider.id, model.id);
    setIsOpen(false);
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 px-4 py-2 bg-gray-800/50 hover:bg-gray-700/50 border border-gray-700 rounded-lg transition-colors backdrop-blur-sm"
      >
        {getProviderIcon(settings.selectedProvider)}
        <span className="text-sm font-medium text-gray-200">
          {currentModel?.name || "Select Model"}
        </span>
        <ChevronDown
          className={`w-4 h-4 text-gray-400 transition-transform ${isOpen ? "rotate-180" : ""}`}
        />
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-2 w-80 bg-gray-800/95 backdrop-blur-sm border border-gray-700 rounded-lg shadow-xl z-50 max-h-96 overflow-y-auto">
          {settings.providers.map((provider) => (
            <div key={provider.id} className="p-2">
              <div className="flex items-center space-x-2 px-2 py-1 text-xs font-semibold text-gray-400 uppercase tracking-wide">
                {getProviderIcon(provider.id)}
                <span>{provider.name}</span>
                {provider.requiresApiKey && !provider.apiKey && (
                  <span className="text-red-400 text-xs">
                    (API Key Required)
                  </span>
                )}
              </div>
              <div className="space-y-1">
                {provider.models.map((model) => (
                  <button
                    key={model.id}
                    onClick={() => handleModelSelect(provider, model)}
                    disabled={provider.requiresApiKey && !provider.apiKey}
                    className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                      settings.selectedModel === model.id &&
                      settings.selectedProvider === provider.id
                        ? "bg-blue-600 text-white"
                        : provider.requiresApiKey && !provider.apiKey
                          ? "text-gray-500 cursor-not-allowed"
                          : "text-gray-300 hover:bg-gray-700/50"
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-medium">{model.name}</span>
                      {model.contextLength && (
                        <span className="text-xs text-gray-500">
                          {model.contextLength.toLocaleString()} ctx
                        </span>
                      )}
                    </div>
                    {model.description && (
                      <p className="text-xs text-gray-500 mt-1">
                        {model.description}
                      </p>
                    )}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};


import React, { useEffect } from "react";
import { Zap } from "lucide-react";

interface CallbackHandlerProps {
  onHandleCallback: () => void;
  isLoading: boolean;
  error: string | null;
}

export const CallbackHandler: React.FC<CallbackHandlerProps> = ({
  onHandleCallback,
  isLoading,
  error,
}) => {
  useEffect(() => {
    onHandleCallback();
  }, [onHandleCallback]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center p-4">
      <div className="text-center">
        <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-2xl">
          <Zap className="w-8 h-8 text-white" />
        </div>

        {isLoading && (
          <>
            <div className="w-8 h-8 border-2 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-white mb-2">
              Authenticating...
            </h2>
            <p className="text-gray-400">
              Please wait while we complete your sign-in
            </p>
          </>
        )}

        {error && (
          <>
            <div className="w-12 h-12 bg-red-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-red-400 text-2xl">×</span>
            </div>
            <h2 className="text-xl font-semibold text-white mb-2">
              Authentication Failed
            </h2>
            <p className="text-red-400 mb-4">{error}</p>
            <button
              onClick={() => (window.location.href = "/")}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
            >
              Return to Login
            </button>
          </>
        )}
      </div>
    </div>
  );
};


import { Message, ModelProvider } from '../types';

export class APIError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'APIError';
  }
}

export const sendMessage = async (
  message: string,
  messages: Message[],
  provider: ModelProvider,
  model: string,
  onChunk?: (chunk: string) => void
): Promise<string> => {
  const { id: providerId, apiKey, baseUrl } = provider;

  try {
    if (providerId === 'openai') {
      return await sendOpenAIMessage(message, messages, model, apiKey!, onChunk);
    } else if (providerId === 'anthropic') {
      return await sendAnthropicMessage(message, messages, model, apiKey!, onChunk);
    } else if (providerId === 'ollama') {
      return await sendOllamaMessage(message, messages, model, baseUrl!, onChunk);
    }
    throw new APIError('Unsupported provider');
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    throw new APIError('Failed to send message: ' + String(error));
  }
};

const sendOpenAIMessage = async (
  message: string,
  messages: Message[],
  model: string,
  apiKey: string,
  onChunk?: (chunk: string) => void
): Promise<string> => {
  const response = await fetch('https://api.openai.com/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model,
      messages: [
        ...messages.map(msg => ({
          role: msg.role,
          content: msg.content,
        })),
        { role: 'user', content: message },
      ],
      stream: !!onChunk,
    }),
  });

  if (!response.ok) {
    throw new APIError(`OpenAI API error: ${response.statusText}`, response.status);
  }

  if (onChunk) {
    return await handleStreamingResponse(response, onChunk);
  }

  const data = await response.json();
  return data.choices[0]?.message?.content || '';
};

const sendAnthropicMessage = async (
  message: string,
  messages: Message[],
  model: string,
  apiKey: string,
  onChunk?: (chunk: string) => void
): Promise<string> => {
  const response = await fetch('https://api.anthropic.com/v1/messages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'x-api-key': apiKey,
      'anthropic-version': '2023-06-01',
    },
    body: JSON.stringify({
      model,
      max_tokens: 4096,
      messages: [
        ...messages.map(msg => ({
          role: msg.role === 'assistant' ? 'assistant' : 'user',
          content: msg.content,
        })),
        { role: 'user', content: message },
      ],
      stream: !!onChunk,
    }),
  });

  if (!response.ok) {
    throw new APIError(`Anthropic API error: ${response.statusText}`, response.status);
  }

  if (onChunk) {
    return await handleStreamingResponse(response, onChunk, 'anthropic');
  }

  const data = await response.json();
  return data.content[0]?.text || '';
};

const sendOllamaMessage = async (
  message: string,
  messages: Message[],
  model: string,
  baseUrl: string,
  onChunk?: (chunk: string) => void
): Promise<string> => {
  const response = await fetch(`${baseUrl}/api/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model,
      messages: [
        ...messages.map(msg => ({
          role: msg.role,
          content: msg.content,
        })),
        { role: 'user', content: message },
      ],
      stream: !!onChunk,
    }),
  });

  if (!response.ok) {
    throw new APIError(`Ollama API error: ${response.statusText}`, response.status);
  }

  if (onChunk) {
    return await handleStreamingResponse(response, onChunk, 'ollama');
  }

  const data = await response.json();
  return data.message?.content || '';
};

const handleStreamingResponse = async (
  response: Response,
  onChunk: (chunk: string) => void,
  provider: 'openai' | 'anthropic' | 'ollama' = 'openai'
): Promise<string> => {
  const reader = response.body?.getReader();
  if (!reader) throw new APIError('No response body');

  const decoder = new TextDecoder();
  let fullResponse = '';

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value, { stream: true });
      const lines = chunk.split('\n').filter(line => line.trim());

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6);
          if (data === '[DONE]') continue;

          try {
            const parsed = JSON.parse(data);
            let content = '';

            if (provider === 'openai') {
              content = parsed.choices?.[0]?.delta?.content || '';
            } else if (provider === 'anthropic') {
              content = parsed.delta?.text || '';
            } else if (provider === 'ollama') {
              content = parsed.message?.content || '';
            }

            if (content) {
              fullResponse += content;
              onChunk(content);
            }
          } catch (e) {
            // Ignore parsing errors for individual chunks
          }
        }
      }
    }
  } finally {
    reader.releaseLock();
  }

  return fullResponse;
};
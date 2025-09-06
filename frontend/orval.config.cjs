module.exports = {
  auth: {
    output: {
      mode: 'tags-split',
      target: 'src/lib/api/auth.ts',
      schemas: 'src/lib/api/models',
      client: 'svelte-query',
      mock: false,
      baseUrl: 'http://localhost:8080',
      override: {
        query: {
          useQuery: true,
          useInfinite: false,
          options: {
            staleTime: 10000,
          },
        },
        mutator: {
          path: './src/lib/api/client.ts',
          name: 'apiClient',
        },
      },
    },
    input: {
      target: '../shared/api/app/v1/app.v1.swagger.yml',
    },
  },
};

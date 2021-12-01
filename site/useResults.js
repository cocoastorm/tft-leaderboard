import useSWR from 'swr';

const fetcher = (...args) => fetch(...args).then(res => res.json())

export default function useResults() {
  const { data, error } = useSWR('/api/leaderboard', fetcher)

  return {
    data,
    isLoading: !error && !data,
    isError: error,
  }
}

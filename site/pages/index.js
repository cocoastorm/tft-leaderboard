import Head from 'next/head';
import useResults from '../src/useResults';
import Participant from '../components/participant';

function Spinner() {
  return (
    <div className="flex flex-col items-center justify-center w-full flex-1 px-20">
      Loading&hellip;
    </div>
  )
}

function Warning({ error }) {
  return (
    <div className="flex flex-col items-center justify-center w-full flex-1 px-20">
      Error: { JSON.stringify(error) }
    </div>
  )
}

function Main() {
  const { data, isLoading, isError } = useResults()

  if (isLoading) {
    return <Spinner />
  }

  if (isError) {
    return <Warning error={ isError } />
  }

  const participants = data.map(
    ({ summoner, rank }, idx) => <Participant key={summoner.id} summoner={summoner} rank={rank} placement={idx+1} />
  );

  return (
    <main className="flex flex-col items-center justify-center flex-1 w-2/3 sm:w-full px-10 py-7">
      { participants }
    </main>
  )
}

export default function Home() {
  return (
    <div className="bg-sky-900 text-white flex flex-col items-center min-h-screen px-2 py-2">
      <Head>
        <title>tft leaderboard</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <header className="pt-2 pb-5">
        <h1 className="text-4xl font-semibold font-sans">tft leaderboard</h1>
        <p className="text-2xl text-center mt-2 font-light font-sans">race to diamond</p>
      </header>

      <Main />
    </div>
  )
}

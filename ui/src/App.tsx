import { useEffect, useState } from 'react'
import { Tweet } from 'react-tweet'

type LikedTweet = {
  time: string
  text: string
  link: string
  createdAt: string
}

const fetchTweetIds = async () => {
  const res = await fetch('https://twi-fav-api-823271554794.asia-northeast1.run.app/liked-tweets')
  if (!res.ok) {
    throw new Error('Failed to fetch liked tweets')
  }
  const tweets: LikedTweet[] = await res.json()
  console.log(tweets)
  return tweets.map((tweet) => tweet.link.split('/').pop()).filter((id) => id !== undefined)
}

function App() {
  const [tweetIds, setTweetIds] = useState<string[]>([])

  useEffect(() => {
    fetchTweetIds().then(setTweetIds).catch(console.error)
  }, [])

  return (
    <>
      {tweetIds.map((id, i) => <Tweet key={i} id={id} />)}
    </>
  )
}

export default App

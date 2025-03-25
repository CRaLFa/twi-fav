import InfiniteScroll from '@rorygudka/react-infinite-scroller'
import { useState } from 'react'
import { Tweet } from 'react-tweet'
import './App.css'
import NotFound from './components/NotFound'
import ScrollTopButton from './components/ScrollTopButton'

type LikedTweet = {
  time: string
  text: string
  link: string
  createdAt: string
}

const API_URL = 'https://twi-fav-api-823271554794.asia-northeast1.run.app'
const LIMIT = 10

const fetchTweetLinks = async (earliestTime: string) => {
  const res = await fetch(`${API_URL}/liked-tweets?earliestTime=${encodeURIComponent(earliestTime)}&limit=${LIMIT}`)
  if (!res.ok) {
    throw new Error('Failed to fetch liked tweets')
  }
  const tweets: LikedTweet[] = await res.json()
  return {
    newTweetLinks: tweets.map((tweet) => tweet.link).filter((link) => link),
    newEarliestTime: tweets.at(-1)?.time,
  }
}

function App() {
  const [tweetLinks, setTweetLinks] = useState<string[]>([])
  const [hasMore, setHasMore] = useState(true)
  const [earliestTime, setEarliestTime] = useState('')

  const loadMore = async () => {
    try {
      const { newTweetLinks, newEarliestTime } = await fetchTweetLinks(earliestTime)
      if (newTweetLinks.length < 1 || !newEarliestTime) {
        setHasMore(false)
        return
      }
      setTweetLinks([...tweetLinks, ...newTweetLinks])
      setEarliestTime(newEarliestTime)
    } catch (e) {
      console.error(e)
      window.alert('ツイートの取得に失敗しました')
    }
  }

  return (
    <>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} loader={<div key={0}>Loading...</div>}>
        {tweetLinks.map((link) => {
          const id = link.split('/').pop()!
          return <Tweet key={id} id={id} onError={() => link} components={NotFound} />
        })}
      </InfiniteScroll>
      <ScrollTopButton />
    </>
  )
}

export default App

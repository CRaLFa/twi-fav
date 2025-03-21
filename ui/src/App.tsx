import InfiniteScroll from '@rorygudka/react-infinite-scroller'
import { useState } from 'react'
import { Tweet } from 'react-tweet'

type LikedTweet = {
  time: string
  text: string
  link: string
  createdAt: string
}

const GET_URL = 'https://twi-fav-api-823271554794.asia-northeast1.run.app/liked-tweets'
const LIMIT = 10

const fetchTweetIds = async (earliestTime: string) => {
  const res = await fetch(`${GET_URL}?earliestTime=${encodeURIComponent(earliestTime)}&limit=${LIMIT}`)
  if (!res.ok) {
    throw new Error('Failed to fetch liked tweets')
  }
  const tweets: LikedTweet[] = await res.json()
  return {
    newTweetIds: tweets.map((tweet) => tweet.link.split('/').pop()).filter((id) => id !== undefined),
    newEarliestTime: tweets.at(-1)?.time,
  }
}

function App() {
  const [tweetIds, setTweetIds] = useState<string[]>([])
  const [hasMore, setHasMore] = useState(true)
  const [earliestTime, setEarliestTime] = useState('')

  const loadMore = async () => {
    try {
      const { newTweetIds, newEarliestTime } = await fetchTweetIds(earliestTime)
      if (newTweetIds.length < 1 || !newEarliestTime) {
        setHasMore(false)
        return
      }
      setTweetIds([...tweetIds, ...newTweetIds])
      setEarliestTime(newEarliestTime)
    } catch (e) {
      console.error(e)
      window.alert('ツイートの取得に失敗しました。')
    }
  }

  return (
    <InfiniteScroll
      loadMore={loadMore}
      hasMore={hasMore}
      loader={<div key={0}>Loading...</div>}
    >
      {tweetIds.map((id, i) => <Tweet key={i} id={id} />)}
    </InfiniteScroll>
  )
}

export default App

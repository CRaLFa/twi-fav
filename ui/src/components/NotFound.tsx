import type { TwitterComponents } from 'react-tweet'

export const NotFound: TwitterComponents = {
  TweetNotFound: ({ error }) => (
    <>
      <div className='react-tweet-theme _root_98iqw_1'>
        <article className='_article_98iqw_21'>
          <div className='_root_16yxa_1'>
            <h3>センシティブな内容もしくは削除済み</h3>
            <p>
              <a href={`https://x.com/i/status/${error}`} target='_blank'>新しいタブで表示する</a>
            </p>
          </div>
        </article>
      </div>
    </>
  ),
}

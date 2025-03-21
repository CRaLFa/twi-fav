import type { TwitterComponents } from 'react-tweet'

export const NotFound: TwitterComponents = {
  TweetNotFound: ({ error }) => (
    <>
      <div className='react-tweet-theme _root_98iqw_1'>
        <article className='_article_98iqw_21'>
          <div className='_root_16yxa_1'>
            <h3>ツイートを表示できません</h3>
            <p>
              以下の原因が考えられます
            </p>
            <ul>
              <li>非公開</li>
              <li>センシティブな内容を含む</li>
              <li>添付メディアの形式が非対応</li>
              <li>削除済み</li>
            </ul>
            <p>
              <a href={`https://x.com/i/status/${error}`} target='_blank' className='_link_1cutb_4'>
                <span className='_text_1cutb_23'>X で表示する</span>
              </a>
            </p>
          </div>
        </article>
      </div>
    </>
  ),
}

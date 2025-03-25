import { TweetContainer, type TwitterComponents } from 'react-tweet'
import not_found_styles from 'react-tweet-theme/tweet-not-found.module.css'
import replies_styles from 'react-tweet-theme/tweet-replies.module.css'

const NotFound: TwitterComponents = {
  TweetNotFound: ({ error }) => {
    const link = (error as string).replace('twitter.com', 'x.com')
    return (
      <TweetContainer>
        <div className={not_found_styles.root}>
          <h3>ツイートを表示できません</h3>
          <p>以下の原因が考えられます</p>
          <ul>
            <li>非公開</li>
            <li>センシティブな内容を含む</li>
            <li>メディアの形式が非対応</li>
            <li>削除済み</li>
          </ul>
          <p>
            <a href={link} target='_blank' rel='noreferrer' className={replies_styles.link}>
              <span className={replies_styles.text}>X で表示する</span>
            </a>
          </p>
        </div>
      </TweetContainer>
    )
  },
}

export default NotFound

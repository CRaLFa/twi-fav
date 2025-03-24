import imgUrl from '../assets/arrow_6353363.png'

const scrollToTop = () =>
  window.scroll({
    top: 0,
    behavior: 'smooth',
  })

export default () => {
  return <img src={imgUrl} alt='scroll to top' className='scroll-top' onClick={scrollToTop} />
}

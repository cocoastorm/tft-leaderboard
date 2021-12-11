import classNames from 'classnames';

export default function Tag(props) {
  const numClass = classNames({'font-normal': true, 'pr-2': true}, props.className)

  return (
    <div>
      <span className={numClass}>{ props.value }</span>
      <span className="font-semibold text-cyan-200">{ props.title }</span>
    </div>
  )
}

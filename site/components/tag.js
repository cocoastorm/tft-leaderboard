export default function Tag(props) {
  const colorClassName = `text-${props.color}-400`;

  return (
    <div>
      <span className={'font-normal pr-2 ' + colorClassName}>{ props.value }</span>
      <span className="font-semibold text-cyan-400">{ props.title }</span>
    </div>
  )
}

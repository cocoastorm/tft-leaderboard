export default function Tag(props) {
  return (
    <div>
      <span className="font-normal">{ props.value }</span>&nbsp;<span className="font-extrabold">{ props.title }</span>
    </div>
  )
}

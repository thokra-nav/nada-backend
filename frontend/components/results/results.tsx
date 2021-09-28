import SearchResult from './searchresult'
import styled from "styled-components";
import {NavGra60} from "../../styles/constants";

const ResultsBox = styled.div`
    border: 2px solid ${NavGra60};
    flex-grow: 1;
    padding: 15px;
`

const Results = () => {
    return (
        <ResultsBox>
            <SearchResult result={{name: 'foo', description: 'bar'}}/>
            <SearchResult result={{name: 'foo', description: 'bar'}}/>
            <SearchResult result={{name: 'foo', description: 'bar'}}/>
            <SearchResult result={{name: 'foo', description: 'bar'}}/>
        </ResultsBox>
    )
}

export default Results